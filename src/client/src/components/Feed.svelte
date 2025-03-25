<script>
    import { onMount } from "svelte";
    import Section from "./Section.svelte";
    import Filter from "./Filter.svelte";
    import { fetchFeedItems, sortFeedItems } from "$lib";
    import { feedState } from "$lib/state.svelte";

    onMount(async () => {
        let results = await fetchFeedItems();
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
<div id="feed-container" class="h-full grid grid-cols-1 gap-y-8 px-4">
    {#if feedState.feedItems.length > 25}
        <Section heading="Latest" items={feedState.feedLatest} />
    {/if}
    <Section heading="All" items={feedState.feedItems} vertical />
</div>
