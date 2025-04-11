local http = require("socket.http")
local ltn12 = require("ltn12")
local json = require("cjson")
local uuid = require("uuid")

uuid.randomseed(os.time())

local function rng()
    local bytes = {}
    for i = 1, 16 do
        bytes[i] = math.random(0, 255)
    end
    return string.char(table.unpack(bytes))
end

uuid.set_rng(rng)

local framework = {
  _debug = true
}

framework.config = {
  loki_url = "http://localhost:8080",
  backoffice_url = "http://localhost:8081",
  timeout = 60
}

function framework.debug(enabled)
  framework._debug = enabled
end

function framework.generate_trace_id()
  return uuid()
end

function framework.generate_request_id()
  return uuid()
end

function framework.log_debug(...)
  if framework._debug then
    print("[DEBUG]", ...)
  end
end

function framework.request(method, url, headers, body)
  local response_body = {}
  local request_body = nil

  if body then
    request_body = json.encode(body)
    framework.log_debug("Request body:", request_body)
  end

  headers = headers or {}
  headers["Content-Type"] = headers["Content-Type"] or "application/json"

  framework.log_debug(string.format("Making %s request to %s", method, url))
  framework.log_debug("Headers:", json.encode(headers))

  local response, status_code, response_headers = http.request {
    url = url,
    method = method,
    headers = headers,
    source = request_body and ltn12.source.string(request_body) or nil,
    sink = ltn12.sink.table(response_body),
    timeout = framework.config.timeout
  }

  local body_str = table.concat(response_body)
  local response_data = nil

  if body_str and body_str ~= "" then
    pcall(function()
      response_data = json.decode(body_str)
    end)
  end

  framework.log_debug(string.format("Response status: %s", status_code))
  framework.log_debug("Response body:", body_str)

  return {
    status = status_code,
    headers = response_headers,
    body = response_data,
    raw_body = body_str
  }
end

function framework.loki_readiness()
  local trace_id = framework.generate_trace_id()
  local request_id = framework.generate_request_id()

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/ready",
    {
      ["X-Trace-ID"] = trace_id,
      ["X-Request-ID"] = request_id
    }
  )

  if not response.body then
    print("Failed to get loki liveness status")
    return nil
  end

  return response
end

function framework.loki_backoffice_readiness()
  local trace_id = framework.generate_trace_id()
  local request_id = framework.generate_request_id()

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/ready",
    {
      ["X-Trace-ID"] = trace_id,
      ["X-Request-ID"] = request_id
    }
  )

  if not response.body then
    print("Failed to get loki-backoffice liveness status")
    return nil
  end

  return response
end

function framework.start_smart_id_auth(country, personal_code)
  local trace_id = framework.generate_trace_id()
  print("Using trace ID: " .. trace_id)

  local response = framework.request(
    "POST",
    framework.config.loki_url .. "/api/auth/smart_id",
    {["X-Trace-ID"] = trace_id},
    {country = country, personal_code = personal_code}
  )

  if not response.body or not response.body.id then
    print("Failed to create session")
    print("Response:", json.encode(response))
    return nil
  end

  print("Session ID: " .. response.body.id)
  return response.body.id, trace_id
end

function framework.check_session_status(session_id, trace_id)
  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/sessions/" .. session_id,
    {["X-Trace-ID"] = trace_id}
  )

  if not response.body then
    print("Failed to get session status")
    return nil
  end

  return response.body.status
end

function framework.wait_for_session_completion(session_id, trace_id, max_attempts)
  max_attempts = max_attempts or 10

  for i = 1, max_attempts do
    local status = framework.check_session_status(session_id, trace_id)
    print("Session status: " .. (status or "unknown"))

    if status == "SUCCESS" then
      return true
    end

    if i == max_attempts then
      print("Timed out waiting for session completion")
      return false
    end

    os.execute("sleep 5")
  end

  return false
end

function framework.complete_auth(session_id, trace_id)
  local response = framework.request(
    "POST",
    framework.config.loki_url .. "/api/sessions/" .. session_id,
    {["X-Trace-ID"] = trace_id}
  )

  if not response.body or not response.body.access_token then
    print("Failed to get access token")
    print("Response:", json.encode(response))
    return nil
  end

  print("Access token received")
  return response.body.access_token
end

function framework.call_backoffice_api(endpoint, token)
  local trace_id = framework.generate_trace_id()
  local request_id = framework.generate_request_id()

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. endpoint,
    {
      ["Authorization"] = "Bearer " .. token,
      ["X-Trace-ID"] = trace_id,
      ["X-Request-ID"] = request_id
    }
  )

  return response
end

function framework.authenticate_with_smart_id(country, personal_code)
  print(string.format("Authenticating with Smart-ID (Country: %s, Personal Code: %s)", country, personal_code))

  local session_id, trace_id = framework.start_smart_id_auth(country, personal_code)
  if not session_id then
    print("Failed to start Smart-ID authentication")
    return nil
  end

  if not framework.wait_for_session_completion(session_id, trace_id) then
    print("Failed waiting for session completion")
    return nil
  end

  local token = framework.complete_auth(session_id, trace_id)
  if not token then
    print("Failed to complete authentication")
    return nil
  end

  print("Authentication successful")
  return token
end

framework.assert = {}

function framework.assert.equals(expected, actual, message)
  if expected ~= actual then
    error(string.format("%s: Expected %s but got %s",
      message or "Assertion failed",
      tostring(expected),
      tostring(actual)
    ))
  end
  return true
end

function framework.assert.not_equals(expected, actual, message)
  if expected == actual then
    error(string.format("%s: Expected value to not equal %s",
      message or "Assertion failed",
      tostring(expected)
    ))
  end
  return true
end

function framework.assert.contains(haystack, needle, message)
  if type(haystack) ~= "string" or type(needle) ~= "string" or not string.find(haystack, needle, 1, true) then
    error(string.format("%s: Expected '%s' to contain '%s'",
      message or "Assertion failed",
      tostring(haystack),
      tostring(needle)
    ))
  end
  return true
end

function framework.assert.status_code(response, expected_status)
  return framework.assert.equals(
    expected_status,
    response.status,
    "Unexpected status code"
  )
end

function framework.assert.has_property(obj, property, message)
  if type(obj) ~= "table" or obj[property] == nil then
    error(string.format("%s: Expected object to have property '%s'",
      message or "Assertion failed",
      tostring(property)
    ))
  end
  return true
end

return framework
