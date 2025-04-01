<script>
    import { fetchFeedItems, sortFeedItems } from "$lib";
    import { feedState } from "$lib/state.svelte";
    import { onMount } from "svelte";
    import Feed from "../components/Feed.svelte";
    import Filter from "../components/Filter.svelte";
    async function updateFeed() {
        let results = await fetchFeedItems(feedState.category);
        console.log("category changed", feedState.category);
        console.log(results);
        feedState.feedItems = sortFeedItems(results.items ? results.items : []);
        console.log("coll", results.collections);
        feedState.feedCollections = results.collections;

        feedState.feedLatest = sortFeedItems(
            results.latest ? results.latest : [],
        );
    }
    $effect(async () => {
        await updateFeed();
    });
    onMount(async () => {
        await updateFeed();
    });
</script>

<Filter />
<Feed />
