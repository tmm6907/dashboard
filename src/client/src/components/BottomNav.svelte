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
    class="fixed bottom-0 left-0 bg-base-200 text-base-content grid grid-cols-5 gap-4 w-full"
>
    <button class="btn btn-ghost py-8 text-ren">
        <div style="display: block;">
            <div><i class="fa-solid fa-bookmark"></i></div>
            <span class="dock-label text-xs">Saved</span>
        </div>
    </button>

    <button class="btn btn-ghost py-8">
        <div style="display: block;">
            <div><i class="fa-solid fa-globe"></i></div>
            <span class="dock-label text-xs">Discover</span>
        </div>
    </button>

    <button id="showBoxBtn" class="btn btn-ghost py-8">
        <div style="display: block;">
            <div><i class="fa-solid fa-plus"></i></div>
            <span class="dock-label text-xs">Add</span>
        </div>
    </button>

    <a href="/" class="btn btn-ghost py-8">
        <div style="display: block;">
            <div><i class="fa-solid fa-house"></i></div>
            <span class="dock-label text-xs">Home</span>
        </div>
    </a>
    <button class="btn btn-ghost py-8" on:click={logout}>
        <div style="display: block;">
            <div><i class="fa-solid fa-user"></i></div>
            <span class="dock-label text-xs">Profile</span>
        </div>
    </button>
</div>
