<script>
    import { getTimeAgo } from "$lib";
    import { selectedItem } from "$lib/state.svelte";
    import { onMount } from "svelte";

    let item = $state({});

    const fetchFeedItem = async () => {
        console.log(selectedItem.id);
        const decodeHTML = (html) => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, "text/html");
            return doc.documentElement.textContent;
        };
        try {
            let id = sessionStorage.getItem("selectedItem");
            if (id === null) {
                console.error("id null");
                return;
            }
            let response = await fetch(
                "https://mashboard.app:8080/api/feeds/items/" + id,
                {
                    credentials: "include",
                },
            );
            if (response.status == 302) {
                console.log("redirect");
                window.location.href = await response.text();
                return;
            }
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            item = await response.json();
            item.description = decodeHTML(item.description);
            console.log(item);
        } catch (e) {
            console.error(e);
        }
    };

    const bookmarkFeedItem = async () => {
        console.log("entered");
        if (!item.id) {
            console.error("feed item id not found");
            return;
        }
        try {
            let response = await fetch(
                `https://mashboard.app:8080/api/feeds/items/${item.id}/bookmark`,
                {
                    method: "POST",
                    credentials: "include",
                },
            );
            if (response.status == 302) {
                window.location.href = await response.text();
                return;
            }
            if (response.status != 200) {
                throw new Error(await response.text());
            }
            item.bookmarked = !item.bookmarked;
        } catch (e) {
            console.error(e);
        }
    };
    $effect(async () => {
        console.log("selected id changed", selectedItem.id);
    });
    onMount(async () => {
        await fetchFeedItem();
    });
</script>

<div class="flex flex-col space-y-12 flex-grow min-h-[75dvh] mb-24 p-4">
    <div class="flex flex-col space-y-4">
        <h1 class="text-xl font-bold">{item.title}</h1>
        <div class="prose flex justify-between gap-1 overflow-y-auto">
            <div class="flex gap-2">
                <div class="max-w-[24ch] text-secondary truncate">
                    {item.feed_name}
                </div>
                |
                <span>{getTimeAgo(item.pub_date ? item.pub_date : "")}</span>
            </div>
            <div class="flex gap-2">
                <button class="btn btn-sm btn-ghost"
                    ><i class="fa fa-plus"></i></button
                >
                <button class="btn btn-sm btn-ghost" onclick={bookmarkFeedItem}
                    ><i
                        class="fa fa-bookmark {item.bookmarked
                            ? 'text-yellow-300'
                            : ''}"
                    ></i></button
                >
            </div>
        </div>
        <img
            src={item.image
                ? item.image
                : "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
            alt={item.altText ? item.altText : "feed item"}
            class="aspect-[3/2]"
        />
        <p class="prose text-base-content text-wrap overflow-y-auto">
            {item.description}
        </p>

        <span class="w-fit">
            <a href={item.link} target="_blank" class="link text-accent">
                Read More...
            </a>
        </span>
    </div>
    <div>
        <div class="grid grid-cols-2 gap-4 mt-auto">
            <button class="btn btn-accent text-neutral"
                ><i class="fa fa-share"></i> Share</button
            >
            <a href="/" class="btn btn-primary">Return to feed</a>
        </div>
    </div>
</div>
