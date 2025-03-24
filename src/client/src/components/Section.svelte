<script>
    import Card from "./Card.svelte";
    let { heading, items, max = 10, vertical = false } = $props();
</script>

<div class="pt-2 pb-8 border-b border-neutral-500">
    <div class="flex justify-between items-center">
        <h2 class="prose prose-2xl">{heading}</h2>
        {#if items.length > 10 && !vertical}
            <span class="bold text-accent"
                >Show All <i class="fa-solid fa-caret-right"></i></span
            >
        {/if}
    </div>
    {#if !vertical}
        <div class="flex gap-8 overflow-x-auto whitespace-nowrap">
            {#each items.slice(0, max) as item}
                <Card podcast={item.name} item={item.item} />
            {/each}
        </div>
    {:else}
        {#each items as item}
            <Card podcast={item.name} item={item.item} />
        {/each}
    {/if}
</div>
