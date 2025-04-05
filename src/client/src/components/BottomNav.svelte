<script>
    import { triggerAlert } from "$lib";
    import { profileData } from "$lib/state.svelte";
    import { onMount } from "svelte";

    const copyToClipboard = async () => {
        navigator.clipboard
            .writeText(profileData.mashboardEmail)
            .then(() => {
                triggerAlert("Copied to clipboard!", { closable: true });
            })
            .catch((err) => {
                // Handle errors (e.g., permission issues)
                console.error("Failed to copy text: " + err);
            });
    };

    const logout = () => {
        fetch("https://api.mashboard.app/auth/logout/", {
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

    const closeSubscribeMenu = () => {
        let subForm = document.getElementById("subscribe-form");
        if (!subForm) return;

        subForm.parentElement?.classList.add("pointer-events-none");
        subForm.classList.remove("slide-in");
        subForm.classList.add("slide-out");
    };

    const openSubscribeMenu = () => {
        let subForm = document.getElementById("subscribe-form");
        if (!subForm) return;
        if (subForm.classList.contains("slide-out")) {
            // Show the subForm
            subForm.parentElement?.classList.remove("pointer-events-none");
            subForm.classList.remove("slide-out");
            subForm.classList.add("slide-in");

            const handleOutsideClick = (e) => {
                if (!subForm.contains(e.target)) {
                    closeSubscribeMenu();
                    document.removeEventListener("click", handleOutsideClick);
                }
            };
            setTimeout(() => {
                document.addEventListener("click", handleOutsideClick);
            }, 0);
        } else {
            // Hide the subForm?
            closeSubscribeMenu();
        }
    };
    onMount(() => {});
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

    <button
        id="showBoxBtn"
        class="btn btn-ghost py-8"
        onclick={openSubscribeMenu}
    >
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
    <div class="relative btn btn-ghost py-8" onclick={toggleMenu}>
        <div style="display: block;">
            <div><i class="fa-solid fa-user"></i></div>
            <span class="dock-label text-xs">Profile</span>
        </div>
        <div
            id="profile-menu"
            class="absolute bg-base-200 -top-60 -left-30 hidden rounded-lg py-4"
        >
            <div></div>
            <ul class="menu pl-0 rounded-box w-[24ch]">
                <li>
                    <div class="dropdown">
                        <div
                            tabindex="0"
                            role="button"
                            onclick={(e) => e.stopPropagation()}
                            class="btn"
                        >
                            Theme
                            <svg
                                width="12px"
                                height="12px"
                                class="inline-block h-2 w-2 fill-current opacity-60"
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 2048 2048"
                            >
                                <path
                                    d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"
                                ></path>
                            </svg>
                        </div>
                        <ul
                            tabindex="0"
                            class="dropdown-content bg-base-300 rounded-box z-1 w-52 p-2 shadow-2xl"
                        >
                            <li>
                                <input
                                    type="radio"
                                    name="theme-dropdown"
                                    class="theme-controller w-full btn btn-sm btn-block btn-ghost justify-start"
                                    aria-label="Default"
                                    value="light"
                                />
                            </li>
                            <li>
                                <input
                                    type="radio"
                                    name="theme-dropdown"
                                    class="theme-controller w-full btn btn-sm btn-block btn-ghost justify-start"
                                    aria-label="Retro"
                                    value="dark"
                                />
                            </li>
                            <li>
                                <input
                                    type="radio"
                                    name="theme-dropdown"
                                    class="theme-controller w-full btn btn-sm btn-block btn-ghost justify-start"
                                    aria-label="Cyberpunk"
                                    value="emerald"
                                />
                            </li>
                        </ul>
                    </div>
                </li>
                <li>
                    <button
                        onclick={copyToClipboard}
                        class="btn btn-ghost text-accent"
                        ><i class="fa-regular fa-clipboard"></i> Mashboard Email</button
                    >
                </li>
                <li><a>View Profile</a></li>
                <li><a>Collections</a></li>
                <li><a>Settings</a></li>
                <li><a>Logout</a></li>
            </ul>
        </div>
    </div>
</div>
