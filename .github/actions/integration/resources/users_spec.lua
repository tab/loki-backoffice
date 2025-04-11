local framework = require("framework")
local auth = require("auth")
local json = require("cjson")

local suite = {}

local test_data = {
  admin_token = nil,
  manager_token = nil,
  user_token = nil,
  user = nil
}

function suite.setup()
  print("Setting up users tests...")
  test_data.admin_token = auth.get_admin_token()
  test_data.manager_token = auth.get_manager_token()
  test_data.user_token = auth.get_user_token()

  return test_data.admin_token ~= nil
end

function suite.test_list()
  print("Test: List users")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/users",
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

function suite.test_create()
  print("Test: Create user")

  local random_number = math.random(10000000000, 99999999999)
  local random_personal_code = tostring(random_number)

  local user = {
    identity_number = "PNOEE-" .. random_personal_code,
    personal_code = random_personal_code,
    first_name = "Luke",
    last_name = "Skywalker",
  }

  local response = framework.request(
    "POST",
    framework.config.backoffice_url .. "/api/backoffice/users",
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    },
    user
  )

  framework.assert.status_code(response, 201)
  framework.assert.has_property(response.body, "id", "Created user missing ID")
  framework.assert.equals(user.identity_number, response.body.identity_number, "Identity number does not match")
  framework.assert.equals(user.personal_code, response.body.personal_code, "Personal code does not match")
  framework.assert.equals(user.first_name, response.body.first_name, "First name does not match")
  framework.assert.equals(user.last_name, response.body.last_name, "Last name does not match")

  test_data.user = response.body

  return response.body.id ~= nil
end

function suite.test_get()
  print("Test: Get user")

  if not test_data.user.id then
    error("No user ID to test")
  end

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/users/" .. test_data.user.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)
  framework.assert.equals(test_data.user.id, response.body.id, "IDs don't match")

  return true
end

function suite.test_update()
  print("Test: Update user")

  if not test_data.user.id then
    error("No user ID to test")
  end

  local user = {
    identity_number = test_data.user.identity_number,
    personal_code = test_data.user.personal_code,
    first_name = "LUKE",
    last_name = "SKYWALKER",
  }

  local response = framework.request(
    "PUT",
    framework.config.backoffice_url .. "/api/backoffice/users/" .. test_data.user.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    },
    user
  )

  framework.assert.status_code(response, 200)
  framework.assert.equals(user.first_name, response.body.first_name, "First name not updated")
  framework.assert.equals(user.last_name, response.body.last_name, "Last name not updated")

  return true
end

function suite.test_delete()
  print("Test: Delete user")

  if not test_data.user.id then
    error("No user ID to test")
  end

  local response = framework.request(
    "DELETE",
    framework.config.backoffice_url .. "/api/backoffice/users/" .. test_data.user.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 204)

  local get_response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/users/" .. test_data.user.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(get_response, 404)

  return true
end

function suite.test_forbidden()
  print("Test: Forbidden access")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/users",
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
    framework.config.backoffice_url .. "/api/backoffice/users",
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
    suite.test_create,
    suite.test_get,
    suite.test_update,
    suite.test_delete,
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
