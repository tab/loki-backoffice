local framework = require("framework")
local auth = require("auth")
local json = require("cjson")

local suite = {}

local test_data = {
  admin_token = nil,
  manager_token = nil,
  user_token = nil
}

function suite.setup()
  print("Setting up tokens tests...")
  test_data.admin_token = auth.get_admin_token()
  test_data.manager_token = auth.get_manager_token()
  test_data.user_token = auth.get_user_token()

  return test_data.admin_token ~= nil
end

function suite.test_list()
  print("Test: List tokens")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/tokens",
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)
  framework.assert.has_property(response.body, "data", "Response missing data array")
  framework.assert.has_property(response.body, "meta", "Response missing pagination metadata")

  return true
end

function suite.test_forbidden()
  print("Test: Forbidden access")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/tokens",
    {
      ["Authorization"] = "Bearer " .. auth.get_user_token(),
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 403)

  return true
end

function suite.test_unauthorized()
  print("Test: Unauthorized access")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/tokens",
    {
      ["Authorization"] = "Bearer " .. auth.get_invalid_token(),
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 401)

  return true
end

function suite.run()
  if not suite.setup() then
    print("❌ Setup failed")
    return false
  end

  local tests = {
    suite.test_list,
    suite.test_forbidden,
    suite.test_unauthorized
  }

  local success = true

  for i, test in ipairs(tests) do
    local test_success, result = pcall(test)

    if not test_success or not result then
      print("❌ Test failed: " .. debug.traceback())
      success = false
      break
    else
      print("✅ Test passed")
    end
  end

  return success
end

return suite
