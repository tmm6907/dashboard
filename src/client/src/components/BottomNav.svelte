<script>
    import { profileData } from "$lib/state.svelte";
    import { onMount } from "svelte";

    const logout = () => {
        fetch("http://localhost:8080/auth/logout/", {
            method: "GET",
            credentials: "include", // Important for cookies/sessions
        })
            .then(() => console.log("Logged out"))
            .catch(console.error);
    };
    const toggleMenu = () => {
        const menu = document.getElementById("profile-menu");
        if (menu.classList.contains("hidden")) {
        }
        menu.classList.toggle("hidden");
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
    <a href="/saved" class="btn btn-ghost py-8 text-ren">
        <div style="display: block;">
            <div><i class="fa-solid fa-bookmark"></i></div>
            <span class="dock-label text-xs">Saved</span>
        </div>
    </a>

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
    <button class="relative btn btn-ghost py-8" on:click={toggleMenu}>
        <div style="display: block;">
            <div><i class="fa-solid fa-user"></i></div>
            <span class="dock-label text-xs">Profile</span>
        </div>
        <div
            id="profile-menu"
            class="absolute bg-base-200 -top-60 -left-30 hidden rounded-lg py-4"
        >
            <div>
                <label for="mashboard-email">{profileData.mashboardEmail}</label
                >
            </div>
            <ul class=" menu rounded-box w-[24ch]">
                <li><a>View Profile</a></li>
                <li><a>Collections</a></li>
                <li><a>Settings</a></li>
                <li><a>Logout</a></li>
            </ul>
        </div>
    </button>
</div>
