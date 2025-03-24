<script>
    import { onMount } from "svelte";

    onMount(() => {
        document.addEventListener("click", (e) => {
            if (e.target.classList.contains("filter-btn")) {
                let filterBtns = document.querySelectorAll(".filter-btn");
                filterBtns.forEach((btn) => {
                    btn.classList.remove("bg-primary");
                });
                e.target.classList.add("bg-primary");
            }
        });
    });
</script>

<div
    class="flex overflow-x-scroll items-center space-x-2"
    style="padding: 1em 0.5em;"
>
    <select name="category-filter" id="category-filter" class="select w-fit">
        <option value="all" selected>All Categories</option>
        <option value="technology">Technology</option>
        <option value="entertainment">Entertainment</option>
        <option value="science">Science</option>
        <option value="politics">Politics</option>
        <option value="lifestyle">Lifestyle</option>
    </select>
    <button class="filter-btn btn btn-sm bg-primary">All</button>
    <button class="filter-btn btn btn-sm">News</button>
    <button class="filter-btn btn btn-sm">Podcasts</button>
    <button class="filter-btn btn btn-sm">Video</button>
    <button class="filter-btn btn btn-sm">Blogs</button>
    <button class="filter-btn btn btn-sm">Entertainment</button>
    <button class="filter-btn btn btn-sm">Web</button>
</div>
