<script>
    import { onMount } from "svelte";
    import Section from "./Section.svelte";
    import Card from "./Card.svelte";

    let feedItems = $state([]);
    let feedLatest = $state([]);
    let feedSaved = $state([]);
    let feedCollections = $state([]);

    const sortFeedItems = (items) => {
        // sort items on item.pubDate
        return items.sort(
            (a, b) => new Date(b.item.pubDate) - new Date(a.item.pubDate),
        );
    };

    const fetchFeeds = async () => {
        console.log("Fetching feed items");
        try {
            const response = await fetch(
                "http://localhost:8080/api/feeds/items",
                {
                    credentials: "include",
                },
            );
            const { latest, items, saved, collections } = await response.json();
            feedItems = sortFeedItems(items ? items : []);
            feedCollections = sortFeedItems(collections ? collections : []);
            feedLatest = sortFeedItems(latest ? latest : []);
            feedSaved = sortFeedItems(saved ? saved : []);
        } catch (err) {
            console.error("Error fetching feeds:", err);
        }
    };
    onMount(async () => {
        await fetchFeeds();
    });
</script>

<div id="feed-container" class="h-full grid grid-cols-1 gap-y-8 px-4">
    <Section heading="Latest" items={feedLatest} />
    <Section heading="All" items={feedItems} vertical />
</div>
