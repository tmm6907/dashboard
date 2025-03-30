<script>
    import { getTimeAgo } from "$lib";
    import { selectedItem } from "$lib/state.svelte";
    import { onMount } from "svelte";
    import Fab from "../../../components/Fab.svelte";

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
                "http://localhost:8080/api/feeds/items/" + id,
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
    $effect(async () => {
        console.log("selected id changed", selectedItem.id);
    });
    onMount(async () => {
        await fetchFeedItem();
    });
</script>

<div class="flex flex-col space-y-4 p-4">
    {#if item}
        <h2 class="text-xl font-bold">{item.title}</h2>
        <div class="prose flex gap-1">
            <div class="max-w-[24ch] text-secondary truncate">
                {item.feedName}
            </div>
            |
            <span>{getTimeAgo(item.pub_date ? item.pub_date : "")}</span>
        </div>
        <img
            src={item.image
                ? item.image
                : "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
            alt={item.altText ? item.altText : "feed item"}
            class="aspect-[3/2]"
        />
        <p class="prose text-neutral-content text-balance">
            {item.description}
        </p>
        <a href={item.link} target="_blank" class="link text-accent"
            >Read More...</a
        >
    {/if}
</div>
<Fab {item} />
