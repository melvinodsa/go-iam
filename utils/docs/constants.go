package docs

const intoDescription = `# Documentation for Go IAM APIs
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
**go-iam** is a lightweight, multi-tenant Identity and Access Management (IAM) server built in **Golang**.
It provides robust authentication and fine-grained authorization for modern applications. With support for custom roles, third-party auth providers, and multi-client setups, *go-iam* gives you full control over access management in a scalable and modular way.


![Go IAM Logo](/docs/goiam.png)

## Resources

> ✅ Admin UI: [go-iam-ui](https://github.com/melvinodsa/go-iam-ui)  
> 🐳 Docker Setup: [go-iam-docker](https://github.com/melvinodsa/go-iam-docker)  
> 🔐 Backend: [go-iam](https://github.com/melvinodsa/go-iam)

## How it works

### Multi-Tenant Architecture
*go-iam* is designed to handle multiple projects or applications within a single instance.
Each project can have its own set of users, roles, auth providers, resources, clients, and policies, allowing for efficient resource management and isolation.

### Authentication and Authorization

*go-iam* supports multiple authentication providers like Google, GitHub, and custom OAuth2 providers.
It uses JWT tokens for secure, stateless authentication. The server can be configured to require authentication for specific routes or resources.
Users can be authenticated using auth providers. Multiple applications can share the same auth provider, allowing for a unified login experience across different projects.
You can create an auth provider and link it to multiple clients. Applications use clients to authenticate users.

![Auth process](/docs/auth-backend.svg)

### Resources and Roles

Anything that needs access control is considered a resource in *go-iam*. Resources can be anything from API endpoints to database records.
Resources can be added to roles. And roles can be assigned to users. There policies can be used to define fine-grained access control.

![Auth process](/docs/resources-roles-users.svg)


<details>
<summary>Examples</summary>

You can create a button called **@ui/my-super-app/button** that allows users to create an account in your application.
You can create a resource called **@ui/my-super-app/resource** that represents the button in your application.
Then you can create a role called **@ui/my-super-app/role** that allows users to access the button. Add the resource to the role.
Then you can assign the role to users. Users can then access the button in your application.
</details>`
