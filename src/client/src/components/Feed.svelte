<script>
    import Section from "./Section.svelte";
    import { feedState } from "$lib/state.svelte";
</script>

<div id="feed-container" class="h-full flex flex-col spacing-y-4 px-6">
    {#if feedState.feedItems.length > 25}
        <Section heading="Latest" items={feedState.feedLatest} />
    {/if}
    <Section heading="All" items={feedState.feedItems} vertical />
</div>
