<script>
    import { onMount } from "svelte";

    let feedQuery = $state(""); // Search query
    let feeds = $state([]); // Store fetched feeds
    let loading = $state(false); // Track loading

    async function fetchFeeds() {
        loading = true;
        try {
            const response = await fetch(
                `https://api.mashboard.app/api/feeds/?query=${feedQuery}`,
                {
                    credentials: "include",
                },
            );
            if (response.status == 302) {
                console.log("redirect");
                window.location.href = await response.text();
                return;
            }
            console.log(response.status);
            const data = await response.json();
            console.log("Fetched feeds:", data); // Debug API response
            feeds = Array.isArray(data.feeds) ? data.feeds : []; // Ensure it's an array
        } catch (err) {
            console.error("Error fetching feeds:", err);
        } finally {
            loading = false;
        }
    }

    function gotoNewFeedForm() {
        let subForm = document.getElementById("subscribe-form");
        subForm?.parentElement?.classList.add("pointer-events-none");
        subForm?.classList.remove("slide-in");
        subForm?.classList.add("slide-out");
        window.location.href = "/new-feed";
    }

    // Fetch whenever `feedName` changes
    $effect(() => {
        fetchFeeds();
    });

    onMount(() => {});
</script>

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
                    bind:value={feedQuery}
                />
            </div>
            <div
                id="matched-feeds"
                class="py-4 grid grid-cols-1"
                style="gap: 1em 0;"
            >
                <button
                    id="create-new-feeed"
                    class="feed-row btn btn-primary flex items-center"
                    onclick={gotoNewFeedForm}
                >
                    <div class="">
                        <i class="fa fa-plus"></i>
                    </div>
                    <span>Create New Feed</span>
                </button>
                <div class="h-64 overflow-y-auto">
                    {#each feeds as feed}
                        <div
                            id={feed.feedId}
                            class="feed-row grid grid-cols-4 items-center border-base-300 cursor-pointer"
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
                                <button class="btn btn-sm btn-soft"
                                    >Follow</button
                                >
                            </div>
                        </div>
                    {/each}
                </div>
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
    .feed-row {
        height: 3rem;
    }
</style>
