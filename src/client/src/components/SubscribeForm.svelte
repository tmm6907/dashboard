<script>
    import { onMount } from "svelte";

    let feedName = $state(""); // Search query
    let feeds = $state([]); // Store fetched feeds
    let loading = $state(false); // Track loading

    async function fetchFeeds() {
        console.log("FETCHING...");
        loading = true;
        try {
            const response = await fetch(
                `http://localhost:8080/api/feeds/?query=${feedName}`,
            );
            const data = await response.json();
            console.log("Fetched feeds:", data); // Debug API response
            feeds = Array.isArray(data.feeds) ? data.feeds : []; // Ensure it's an array
        } catch (err) {
            console.error("Error fetching feeds:", err);
        } finally {
            loading = false;
        }
    }

    // Fetch whenever `feedName` changes
    $effect(() => {
        fetchFeeds();
    });

    const openModal = () => {
        const modal = document.getElementById("my_modal_3");
        modal?.showModal();
    };

    onMount(() => {
        document.addEventListener("click", (e) => {
            if (e?.target?.id === "create-new-feed-btn") {
                openModal();
            }
        });
    });
</script>

<dialog id="my_modal_3" class="modal" inert>
    <div class="modal-box">
        <form method="dialog">
            <button
                class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
                ><i class="fa-solid fa-x"></i></button
            >
        </form>
        <h3 class="text-lg font-bold">Hello!</h3>
        <p class="py-4">Press ESC key or click on ✕ button to close</p>
    </div>
</dialog>

<div class="fixed bottom-0 left-0 w-full pointer-events-none">
    <form action="" id="subscribe-form" class="slide-out">
        <div
            class="
            subscribe-form
            gap-y-4 bg-base-100
            h-[100dvh]
            border-t rounded-xl
            "
        >
            <div class="bg-base-300 mb-4">
                <input
                    id="feed-subcribe"
                    type="text"
                    class="border border-base-300 rounded p-2 w-full"
                    placeholder="Search by name, keyword, or url"
                    bind:value={feedName}
                />
            </div>
            <div
                id="matched-feeds"
                class="py-4 grid grid-cols-1 *:h-12"
                style="gap: 1em 0;"
            >
                <div
                    id="create-new-feeed"
                    role="button"
                    class="btn bg-primary"
                    style="gap:0.5em;"
                >
                    <div class="flex items-center">
                        <i class="fa fa-plus"></i>
                    </div>
                    <span>Create New Feed</span>
                </div>
                {#each feeds as feed}
                    <div
                        id={feed.feedId}
                        role="button"
                        onclick={() => console.log("Clicked:", feed)}
                        class="grid grid-cols-4 items-center border-base-300 cursor-pointer"
                        style="gap:0.5em;"
                    >
                        <div class="col-span-1 items-center align-middle">
                            <div class="w-fit m-auto">
                                <img
                                    class="w-[3rem]"
                                    src={feed.image ||
                                        "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
                                    alt="feed item image"
                                />
                            </div>
                        </div>
                        <div class="col-span-2">
                            <span>{feed.title}</span>
                        </div>
                        <div class="flex justify-end pr-2">
                            <button class="btn btn-sm btn-neutral-content"
                                >Follow</button
                            >
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    </form>
</div>

<style>
    .subscribe-form {
        padding: 2em 1em;
        border-top: 2px solid var(--color-base-200);
        border-radius: 2rem;
    }
</style>
