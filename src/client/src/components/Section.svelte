<script>
    import Card from "./Card.svelte";
    let { heading, items, max = 10, vertical = false } = $props();
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
        class="{vertical ? '' : 'flex'} gap-8 overflow-x-auto whitespace-nowrap"
    >
        {#each items.slice(0, vertical ? items.length : max) as item}
            <Card
                podcast={item.name}
                item={item.item}
                small={vertical ? "" : true}
            />
        {/each}
    </div>
</div>
