import { apiFetch, loggedIn } from "./api/api.mjs";

function addButton({ parent, textContent, click }) {
  const button = document.createElement("button");
  button.textContent = textContent;
  button.addEventListener("click", click);
  parent.appendChild(button);
}

document.addEventListener("DOMContentLoaded", async () => {
  const buttonMenu = document.getElementById("button-menu");

  if (await loggedIn()) {
    addButton({
      parent: buttonMenu,
      textContent: "Logout",
    });
    addButton({
      parent: buttonMenu,
      textContent: "Lists",
    });
  } else {
    addButton({
      parent: buttonMenu,
      textContent: "login",
      click: async () => {
        window.location.href = "/login.html";
      },
    });

    addButton({
      parent: buttonMenu,
      textContent: "New List as Guest",
      click: async () => {
        const json = await apiFetch("api/v1/lists", {
          method: "POST",
          body: JSON.stringify({}),
          headers: {
            "Content-Type": "application/json",
          },
        });
        // redirect to the new list page
        window.location.href = `/list.html?ListID=${encodeURIComponent(json.ID)}`;
      },
    });
  }
});
