<script>
    import { getTimeAgoAll } from "$lib";
    import { onMount } from "svelte";

    // @ts-ignore
    let feedData = $state({
        url: "",
        title: "",
        desc: "",
        img: {
            url: "",
            title: "",
        },
        items: [],
        collections: [],
    });
    let selectedCollection = $state("");
    let validFeed = $state(false);

    async function fetchFeedData() {
        try {
            console.log(feedData.url);
            let response = await fetch(
                "http://50.116.53.73:8080/api/feeds/find",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                    body: JSON.stringify({ url: feedData.url }),
                },
            );
            if (response.status == 302) {
                console.log("redirect");
                window.location.href = await response.text();
                return;
            }
            if (response.status != 200) {
                let err = await response.text();
                throw new Error(`unable to create feed: ${err}`);
            }
            validFeed = true;
            let { title, description, image, items } = await response.json();
            feedData.title = title;
            feedData.desc = description;
            feedData.img = image;
            feedData.items = items;
        } catch (e) {
            console.error(e);
        }
    }

    function followFeed() {
        let form = document.getElementById("new-feed-form");
        if (!form) {
            console.error("No form found!");
            return;
        }
        // @ts-ignore
        let formData = new FormData(form);
        if (formData.get("collection") == "other") {
            formData.set("collection", String(formData.get("new-collection")));
            formData.delete("new-collection");
        }
        let body = JSON.stringify(Object.fromEntries(formData));
        console.log("body", body);
        fetch("http://50.116.53.73:8080/api/feeds/follow", {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: body,
        }).then(async (resp) => {
            if (resp.status == 302) {
                console.log("redirect");
                window.location.href = await resp.text();
                return;
            }
            if (resp.status == 200) {
                let my_modal_2 = document.getElementById("my_modal_2");
                my_modal_2.style.display = "none";
                window.location.href = "/";
            }
            let err = await resp.text();
            throw new Error(err);
        });
    }

    function previewFeed() {}

    $effect(async () => {
        if (feedData.url.length >= 3) {
            await fetchFeedData();
        }
    });
    onMount(() => {
        let my_modal_2 = document.getElementById("my_modal_2");
    });
</script>

<dialog id="my_modal_2" class="modal">
    <div class="modal-box">
        <h2 class="text-lg font-bold text-primary">{feedData.title}</h2>
        <div class="h-[50vh] overflow-y-auto">
            <div class="px-4">
                {#if feedData.img}
                    <img
                        class="w-[full] aspect-[3/2]"
                        width="100%"
                        src={feedData.img.url}
                        alt={feedData.img.title}
                    />
                {/if}
            </div>
            <div class="py-4">
                <h3 class="text-secondary">Description</h3>
                <summary class="prose list-none">{feedData.desc}</summary>
            </div>
            <div>
                <h3 class="text-secondary">Latest Posts</h3>
                <ul class="list bg-base-100">
                    {#each feedData.items.slice(0, 10) as item}
                        <li class="grid grid-cols-5 gap-2 py-4">
                            <div>
                                <span class="text-xs prose"
                                    >{getTimeAgoAll(item.published)}</span
                                >
                            </div>
                            <div class="col-span-4">
                                <span>{item.title}</span>
                            </div>
                        </li>
                    {/each}
                </ul>
            </div>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button>Close</button>
    </form>
</dialog>

<form id="new-feed-form">
    <div class="grid grid-cols-1 gap-8 py-2">
        <div>
            <label class="prose" for="new-feed-link"
                >URL<span class="asterisk text-primary ml-1">*</span></label
            >
            <input
                type="text"
                name="link"
                id="new-feed-link"
                class="input bg-base-300"
                required
                bind:value={feedData.url}
            />
            <!-- <div id="feed-image-container"></div> -->
        </div>

        <div>
            <label class="prose" for="new-feed-title"
                >Title<span class="asterisk text-primary ml-1">*</span></label
            >
            <input
                type="text"
                name="title"
                id="new-feed-title"
                class="input bg-base-300"
                required
                bind:value={feedData.title}
            />
        </div>
        <div>
            <label class="prose" for="new-feed-desc">Description</label>
            <textarea
                name="desc"
                id="new-feed-desc"
                class="textarea bg-base-300"
                bind:value={feedData.desc}
            />
        </div>
        <div>
            <label class="prose" for="new-feed-collection">Collection</label>
            <select
                name="collection"
                id="new-feed-collection"
                class="select bg-base-300 mb-4"
                bind:value={selectedCollection}
            >
                {#each feedData.collections as collection}
                    <option value={collection}>{collection}</option>
                {/each}
                <option value="other">Other</option>
            </select>
            <input
                type="text"
                name="new-collection"
                id="new-collection"
                class="input bg-base-300"
                disabled={selectedCollection != "other"}
            />
        </div>

        <div class="grid grid-cols-1 gap-y-2">
            <div class="grid grid-cols-2 gap-4 mt-12">
                <button
                    onclick={window.history.back()}
                    type="button"
                    class="btn">Cancel</button
                >
                <button
                    onclick={my_modal_2.showModal()}
                    type="button"
                    class="btn btn-accent"
                    disabled={!validFeed}>Preview</button
                >
            </div>
            <button class="btn btn-primary" type="button" onclick={followFeed}
                >Create</button
            >
        </div>
    </div>
</form>

<style>
    input,
    select,
    textarea {
        width: 100%;
    }
</style>
