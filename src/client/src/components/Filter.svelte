<script>
    import { onMount } from "svelte";
    import { feedState } from "$lib/state.svelte";
    let filter = $state("All Categories");

    onMount(() => {
        document.addEventListener("click", (e) => {
            let btn = e.target.closest(".filter-btn");
            if (!btn) return;
            console.log("clicked");
            let container = btn.parentElement;
            console.log(container);
            let oldBtn = container.querySelector(".btn-primary");
            oldBtn.classList.remove("btn-primary");
            feedState.category = btn.textContent;
            btn.classList.add("btn-primary");
        });
    });
</script>

<div
    class="flex overflow-x-auto items-center space-x-2"
    style="padding: 1em 0.5em;"
>
    <select
        name="category-filter"
        id="category-filter"
        class="select w-fit"
        bind:value={filter}
    >
        <option value="all" selected>All Categories</option>
        <option value="technology">Technology</option>
        <option value="entertainment">Entertainment</option>
        <option value="science">Science</option>
        <option value="politics">Politics</option>
        <option value="lifestyle">Lifestyle</option>
    </select>
    <button class="filter-btn btn btn-sm btn-primary">All</button>
    <button class="filter-btn btn btn-sm">News</button>
    <button class="filter-btn btn btn-sm">Podcasts</button>
    <button class="filter-btn btn btn-sm">Video</button>
    <button class="filter-btn btn btn-sm">Blogs</button>
    <button class="filter-btn btn btn-sm">Entertainment</button>
    <button class="filter-btn btn btn-sm">Web</button>
</div>
