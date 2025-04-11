local framework = require("framework")
local auth = require("auth")
local json = require("cjson")

local suite = {}

local test_data = {
  admin_token = nil,
  manager_token = nil,
  user_token = nil,
  role = nil
}

function suite.setup()
  print("Setting up roles tests...")
  test_data.admin_token = auth.get_admin_token()
  test_data.manager_token = auth.get_manager_token()
  test_data.user_token = auth.get_user_token()

  return test_data.admin_token ~= nil
end

function suite.test_list()
  print("Test: List roles")

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/roles",
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
  print("Test: Create role")

  local role = {
    name = "developer" .. os.time(),
    description = "Developer test role created by integration test"
  }

  local response = framework.request(
    "POST",
    framework.config.backoffice_url .. "/api/backoffice/roles",
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    },
    role
  )

  framework.assert.status_code(response, 201)
  framework.assert.has_property(response.body, "id", "Created role missing ID")
  framework.assert.equals(role.name, response.body.name, "Name does not match")
  framework.assert.equals(role.description, response.body.description, "Description does not match")

  test_data.role = response.body

  return response.body.id ~= nil
end

function suite.test_get()
  print("Test: Get role")

  if not test_data.role.id then
    error("No role ID to test")
  end

  local response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/roles/" .. test_data.role.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)
  framework.assert.equals(test_data.role.id, response.body.id, "IDs don't match")

  return true
end

function suite.test_update()
  print("Test: Update role")

  if not test_data.role.id then
    error("No role ID to test")
  end

  local role = {
    name = "developer-updated" .. os.time(),
    description = "Developer updated test role"
  }

  local response = framework.request(
    "PUT",
    framework.config.backoffice_url .. "/api/backoffice/roles/" .. test_data.role.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    },
    role
  )

  framework.assert.status_code(response, 200)
  framework.assert.equals(role.name, response.body.name, "Name not updated")
  framework.assert.equals(role.description, response.body.description, "Description not updated")

  return true
end

function suite.test_delete()
  print("Test: Delete role")

  if not test_data.role.id then
    error("No role ID to test")
  end

  local response = framework.request(
    "DELETE",
    framework.config.backoffice_url .. "/api/backoffice/roles/" .. test_data.role.id,
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 204)

  local get_response = framework.request(
    "GET",
    framework.config.backoffice_url .. "/api/backoffice/roles/" .. test_data.role.id,
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
    framework.config.backoffice_url .. "/api/backoffice/roles",
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
    framework.config.backoffice_url .. "/api/backoffice/roles",
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
