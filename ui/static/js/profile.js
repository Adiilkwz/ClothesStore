document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");
    const authSection = document.getElementById("auth-section");
    const profileSection = document.getElementById("profile-section");

    if (!token) {
        authSection.style.display = "flex";
        setupAuthForms();
    } else {
        profileSection.style.display = "block";
        loadUserProfile(token);
    }
});

function setupAuthForms() {
    document.getElementById("login-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const email = document.getElementById("login-email").value;
        const password = document.getElementById("login-password").value;

        await handleAuth("/login", { email, password });
    });

    document.getElementById("signup-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const userData = {
            name: document.getElementById("s-name").value,
            email: document.getElementById("s-email").value,
            password: document.getElementById("s-password").value,
            address: document.getElementById("s-address").value
        };

        const success = await handleAuth("/signup", userData);
        if(success) alert("Account created! Please login.");
    });
}

async function handleAuth(endpoint, data) {
    try {
        const res = await fetch(endpoint, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data)
        });

        if (res.ok) {
            if (endpoint === "/login") {
                const resData = await res.json();
                localStorage.setItem("token", resData.token);
                localStorage.setItem("role", resData.role);
                window.location.reload();
            }
            return true;
        } else {
            alert("Action failed. Check details.");
            return false;
        }
    } catch (err) {
        console.error(err);
        return false;
    }
}

async function loadUserProfile(token) {
    try {
        const res = await fetch("/api/users/me", {
            headers: { "Authorization": "Bearer " + token }
        });

        if (res.ok) {
            const user = await res.json();
            document.getElementById("user-name").innerText = `Hello, ${user.name}`;
            document.getElementById("user-email").innerText = user.email;
            document.getElementById("user-address").innerText = user.address;
        } else {
            logout();
        }
    } catch (err) {
        console.error(err);
    }
}

function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    window.location.reload();
}