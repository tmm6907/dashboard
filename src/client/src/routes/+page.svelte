<script>
    import { fetchFeedItems, sortFeedItems } from "$lib";
    import { feedState } from "$lib/state.svelte";
    import Feed from "../components/Feed.svelte";
    import Filter from "../components/Filter.svelte";
    $effect(async () => {
        let results = await fetchFeedItems(feedState.category);
        console.log("category changed", feedState.category);
        feedState.feedItems = sortFeedItems(results.items ? results.items : []);
        feedState.feedCollections = sortFeedItems(
            results.collections ? results.collections : [],
        );
        feedState.feedLatest = sortFeedItems(
            results.latest ? results.latest : [],
        );
        feedState.feedSaved = sortFeedItems(results.saved ? results.saved : []);
    });
</script>

<Filter />
<Feed />
