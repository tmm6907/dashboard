<script>
    import Card from "./Card.svelte";
    let { heading, items, max = 10, vertical = false, feed = true } = $props();
</script>

<div class="pt-2 pb-4 border-b border-neutral-500">
    <div class="flex justify-between items-center bold">
        <h2 class="prose prose-2xl">{heading}</h2>
        {#if items.length > max && !vertical}
            <span class="text-accent"
                >Show All <i class="fa-solid fa-caret-right"></i></span
            >
        {/if}
    </div>
    <div
        class="{vertical
            ? ''
            : 'flex'} section-container gap-8 overflow-x-auto pb-8 whitespace-nowrap"
    >
        {#if feed}
            {#each items.slice(0, vertical ? items.length : max) as item}
                <Card {item} />
            {/each}
        {:else}
            {#each items.slice(0, vertical ? items.length : max) as item}
                <Card {item} {feed} />
            {/each}
            <Card {feed} />
        {/if}
    </div>
</div>

<style>
    .section-container::-webkit-scrollbar {
        display: none; /* Hide scrollbar for Chrome, Safari */
    }
</style>
