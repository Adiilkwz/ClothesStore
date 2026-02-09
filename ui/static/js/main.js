document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("token");
    const role = localStorage.getItem("role");

    if (token) {
        const loginLink = document.getElementById("login-link");
        if (loginLink) loginLink.style.display = "none";

        const logoutLink = document.getElementById("logout-link");
        if (logoutLink) logoutLink.style.display = "block";

        const profileLink = document.getElementById("profile-link");
        if (profileLink) profileLink.style.display = "block";

        if (role === "admin") {
            const adminLink = document.getElementById("admin-link");
            if (adminLink) adminLink.style.display = "block";
        }
    }
});

function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role"); // 👈 Clear role on logout
    window.location.href = "/";
}