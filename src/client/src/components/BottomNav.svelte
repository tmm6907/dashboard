<script>
    import { onMount } from "svelte";

    const logout = () => {
        fetch("http://localhost:8080/auth/logout/", {
            method: "GET",
            credentials: "include", // Important for cookies/sessions
        })
            .then(() => console.log("Logged out"))
            .catch(console.error);
    };
    onMount(() => {
        document.getElementById("showBoxBtn")?.addEventListener("click", () => {
            let subForm = document.getElementById("subscribe-form");

            if (subForm?.classList.contains("slide-out")) {
                // Show the subForm?
                subForm.parentElement?.classList.remove("pointer-events-none");
                subForm?.classList.remove("slide-out");
                subForm?.classList.add("slide-in");
            } else {
                // Hide the subForm?
                subForm?.parentElement?.classList.add("pointer-events-none");
                subForm?.classList.remove("slide-in");
                subForm?.classList.add("slide-out");
            }
        });
    });
</script>

<div
    class="fixed bg-base-100 text-base-content grid grid-flow-col gap-2 w-full"
    style="padding: 1em 0; gap: 0.5em; bottom: 0; left:0;"
>
    <button class="btn btn-ghost text-ren">
        <div style="display: block;">
            <div><i class="fa-solid fa-bookmark"></i></div>
            <span class="dock-label">Saved</span>
        </div>
    </button>

    <button class="btn btn-ghost">
        <div style="display: block;">
            <div><i class="fa-solid fa-globe"></i></div>
            <span class="dock-label">Discover</span>
        </div>
    </button>

    <button id="showBoxBtn" class="btn btn-ghost">
        <div style="display: block;">
            <div><i class="fa-solid fa-plus"></i></div>
            <span class="dock-label">Add</span>
        </div>
    </button>

    <button class="btn btn-ghost">
        <div style="display: block;">
            <div><i class="fa-solid fa-house"></i></div>
            <span class="dock-label">Home</span>
        </div>
    </button>
    <button class="btn btn-ghost" on:click={logout}>
        <div style="display: block;">
            <div><i class="fa-solid fa-user"></i></div>
            <span class="dock-label">Profile</span>
        </div>
    </button>
</div>
