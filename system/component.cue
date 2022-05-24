package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

component: schema.#component & {
	handle: "system"

	resources: {
		"attachment":            attachment
		"application":           application
		"apigw-route":           apigw_route
		"apigw-filter":          apigw_filter
		"auth-client":           auth_client
		"auth-confirmed-client": auth_confirmed_client
		"auth-session":          auth_session
		"auth-oa2token":         auth_oa2token
		"credential":            credential
		"queue":                 queue
		"queue_message":         queue_message
		"reminder":              reminder
		"report":                report
		"resource-translation":  resource_translation
		"role":                  role
		"role_member":           role_member
		"settings":              settings
		"template":              template
		"user":                  user
		"dal_connection":        dal_connection
	}

	rbac: operations: {
		"action-log.read": description: "Access to action log"

		"settings.read": description:       "Read system settings"
		"settings.manage": description:     "Manage system settings"
		"auth-client.create": description:  "Create auth clients"
		"auth-clients.search": description: "List, search or filter auth clients"

		"role.create": description:  "Create roles"
		"roles.search": description: "List, search or filter roles"

		"user.create": description:  "Create users"
		"users.search": description: "List, search or filter users"

		"dal-connection.create": description:  "Create DAL connections"
		"dal-connections.search": description: "List, search or filter DAL connections"

		"application.create": description:      "Create applications"
		"applications.search": description:     "List, search or filter auth clients"
		"application.flag.self": description:   "Manage private flags for applications"
		"application.flag.global": description: "Manage global flags for applications"

		"template.create": description:  "Create template"
		"templates.search": description: "List, search or filter templates"

		"report.create": description:  "Create report"
		"reports.search": description: "List, search or filter reports"

		"reminder.assign": description: " Assign reminders"

		"queue.create": description:  "Create messagebus queues"
		"queues.search": description: "List, search or filter messagebus queues"

		"apigw-route.create": description:  "Create API gateway route"
		"apigw-routes.search": description: "List search or filter API gateway routes"

		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
