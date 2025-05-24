<script>
    import { triggerAlert } from "$lib";
    import { profileData } from "$lib/state.svelte";
    import { onMount } from "svelte";

    // $effect(() => {
    //     localStorage.setItem("theme", appState.theme);
    // });
    // onMount(() => {
    //     appState.theme = localStorage.getItem("theme") || "dark";
    // });

    let toggled = $state(false);

    const saveTheme = () => {
        localStorage.setItem("isToggled", JSON.stringify(toggled));
    };

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
            .then(() => {
                console.log("Logged out");
                window.location.href = "https://mashboard.app/login";
            })
            .catch((e) => console.error(e));
    };
    const toggleMenu = () => {
        const menu = document.getElementById("profile-menu");
        if (!menu) {
            console.error("profile menu not found");
            return;
        }
        if (menu.classList.contains("hidden")) {
            const handleOutsideClick = (e) => {
                if (!menu.contains(e.target)) {
                    menu.classList.add("hidden");
                    document.removeEventListener("click", handleOutsideClick);
                }
            };
            menu?.classList.remove("hidden");
            setTimeout(() => {
                document.addEventListener("click", handleOutsideClick);
            }, 0);
            return;
        }
        menu.classList.add("hidden");
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
    onMount(() => {
        toggled = JSON.parse(localStorage.getItem("isToggeled") || "false");
    });
</script>

<div
    class="fixed bottom-0 left-0 bg-base-200 text-base-content grid grid-cols-5 gap-4 w-full"
>
    <a href="https://mashboard.app/saved" class="btn btn-ghost py-8 text-ren">
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

    <a href="https://mashboard.app" class="btn btn-ghost py-8">
        <div style="display: block;">
            <div><i class="fa-solid fa-house"></i></div>
            <span class="dock-label text-xs">Home</span>
        </div>
    </a>
    <div class="relative flex items-center">
        <div class="block btn btn-ghost my-auto" onclick={toggleMenu}>
            <div><i class="fa-solid fa-user"></i></div>
            <div class=" dock-label text-xs">Profile</div>
        </div>
        <div
            id="profile-menu"
            class="absolute bg-base-200 -top-70 -left-30 hidden rounded-lg py-4"
        >
            <ul class="menu px-0 rounded-box w-[24ch]">
                <li>
                    <button
                        onclick={copyToClipboard}
                        class="btn btn-ghost text-accent"
                        ><i class="fa-regular fa-clipboard"></i> Mashboard Email</button
                    >
                </li>
                <li class="">
                    <form class="flex justify-between rounded-circle">
                        <div>
                            <label for="theme-selector">Theme</label>
                        </div>
                        <div>
                            <input
                                id="theme-selector"
                                type="checkbox"
                                onchange={saveTheme}
                                value="emerald"
                                checked={toggled}
                                class="toggle theme-controller"
                            />
                        </div>
                    </form>
                </li>
                <li><a>View Profile</a></li>
                <li><a>Collections</a></li>
                <li><a>Settings</a></li>
                <li><button onclick={logout}>Logout</button></li>
            </ul>
        </div>
    </div>
</div>
