import { showError } from "./utils/error.mjs";

document.addEventListener("DOMContentLoaded", async () => {
  const loginButton = document.getElementById("loginButton");
  const nameInput = document.getElementById("nameInput");
  const passwordInput = document.getElementById("passwordInput");

  loginButton.addEventListener("click", async (event) => {
    event.preventDefault();
    const name = nameInput.value;
    const password = passwordInput.value;
    try {
      const response = await fetch("api/v1/session", {
        method: "POST",
        body: JSON.stringify({
          name: name,
          password: password,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        throw new Error(`API Error, ${response.status}`);
      }
    } catch (error) {
      showError(error, "login failed");
    }
    window.location.href = "/index.html";
  });
});
