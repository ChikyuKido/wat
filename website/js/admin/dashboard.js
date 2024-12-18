document.addEventListener("DOMContentLoaded", () => {
    loadUsers();
    loadPermissions();
    loadRoles();
});

function loadUsers() {
    fetch("/api/v1/admin/users/list")
        .then(response => response.json())
        .then(users => {
            const userTable = document.getElementById("userTable");
            userTable.innerHTML = "";
            users.forEach(renderUserRow);
        })
        .catch(error => console.error("Error loading users:", error));
}

function renderUserRow(user) {
    const row = document.createElement("tr");

    row.innerHTML = `
        <td>${user.id}</td>
        <td>${user.email}</td>
        <td>${user.username}</td>
        <td>${user.verified}</td>
        <td>${user.permissions.map(value => `<button class="button mr-1 is-small is-info" onclick="removePermissionFromUser(${user.id},${value.id})">${value.name}</button>`).join("")}</td>
        <td>
            <button class="button is-small is-danger" onclick="deleteUser(${user.id})">Delete</button>
            <button class="button is-small is-primary" onclick="addPermissions(${user.id})">Add Permission</button>
            <button class="button is-small is-primary" onclick="manageRoles(${user.id})">Manage Roles</button>
        </td>
    `;
    document.getElementById("userTable").appendChild(row);
}

function loadPermissions() {
    fetch("/api/v1/admin/permissions/list")
        .then(response => response.json())
        .then(permissions => {
            const select = document.getElementById("permission-dropdown");
            select.innerHTML = "";
            permissions.forEach(renderPermissionOption);
        })
        .catch(error => console.error("Error loading permissions:", error));
}

function renderPermissionOption(permission) {
    const option = document.createElement("option");
    option.value = permission.id;
    option.text = permission.name;
    document.getElementById("permission-dropdown").appendChild(option);
}

function loadRoles() {
    fetch("/api/v1/admin/roles/list")
        .then(response => response.json())
        .then(roles => {
            const select = document.getElementById("role-dropdown");
            select.innerHTML = "";
            roles.forEach(renderRoleOption);
        })
        .catch(error => console.error("Error loading roles:", error));
}

function renderRoleOption(role) {
    const option = document.createElement("option");
    option.value = role.id;
    option.text = role.name;
    document.getElementById("role-dropdown").appendChild(option);
}

function addPermissions(userId) {
    document.getElementById("add-permission").onclick = () => {
        const selectedPermissions = getSelectedPermissions();
        const payload = { user_id: userId, permission_id: selectedPermissions };
        addPermissionToUser(payload);
    };
    document.getElementById("permission-modal").classList.add("is-active");
}

function manageRoles(userId) {
    document.getElementById("add-role").onclick = () => {
        const selectedRoles = getSelectedRoles();
        const overwriteOldPermissions = document.getElementById("overwrite-old-permissions").value === "on";
        const payload = { user_id: userId, role_id: selectedRoles, overwrite_old_permissions: overwriteOldPermissions };
        addRoleToUser(payload);
    };
    document.getElementById("remove-role").onclick = () => {
        const selectedRoles = getSelectedRoles();
        removeRoleFromUser(userId, selectedRoles);
    };
    document.getElementById("role-modal").classList.add("is-active");
}

function getSelectedPermissions() {
    const select = document.getElementById("permission-dropdown");
    return parseInt(select.value);
}

function getSelectedRoles() {
    const select = document.getElementById("role-dropdown");
    return parseInt(select.value);
}

function addPermissionToUser(payload) {
    fetch("/api/v1/admin/users/addPermissionToUser", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    })
        .then(response => response.json())
        .then(data => {
            console.log("Permission added:", data);
            loadUsers();
            document.getElementById("permission-modal").classList.remove("is-active");
        })
        .catch(error => console.error("Error:", error));
}

function addRoleToUser(payload) {
    fetch("/api/v1/admin/users/addRoleToUser", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    })
        .then(response => response.json())
        .then(data => {
            console.log("Role added:", data);
            loadUsers();
            document.getElementById("role-modal").classList.remove("is-active");
        })
        .catch(error => console.error("Error:", error));
}

function deleteUser(userId) {
    if (confirm("Are you sure you want to delete this user?")) {
        fetch(`/api/v1/admin/users/deleteUser`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ user_id: userId })
        })
            .then(response => response.json())
            .then(data => {
                console.log("User deleted:", data);
                loadUsers();
            })
            .catch(error => console.error("Error:", error));
    }
}

function removePermissionFromUser(userId, permissionId) {
    fetch(`/api/v1/admin/users/removePermissionFromUser`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ user_id: userId, permission_id: permissionId })
    })
        .then(response => response.json())
        .then(data => {
            console.log("Permission removed:", data);
            loadUsers();
        })
        .catch(error => console.error("Error:", error));
}

function removeRoleFromUser(userId, roleId) {
    fetch(`/api/v1/admin/users/removeRoleFromUser`, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ user_id: userId, role_id: roleId })
    })
        .then(response => response.json())
        .then(data => {
            console.log("Role removed:", data);
            loadUsers();
        })
        .catch(error => console.error("Error:", error));
}