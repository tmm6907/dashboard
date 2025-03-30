<script>
    import { getTimeAgo } from "$lib";
    let { podcast, item } = $props();
    const goToFeedItem = () => {
        sessionStorage.setItem("selectedItem", `${item.id}`);
        window.location.href = "/items/" + item.id;
    };
</script>

<div
    role="button"
    class="text-base-content pb-4 min-w-[32ch] max-w-[48ch]"
    on:click={goToFeedItem}
>
    <img
        src={item.image
            ? item.image
            : "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
        alt={item.altText ? item.altText : "feed item"}
        loading="lazy"
        class="h-48 w-full object-fill"
    />
    <span class=" text-wrap line-clamp-2">{item.title}</span>
    <div class="flex justify-between py-2">
        <div class="prose">
            <div class="max-w-[24ch] text-secondary truncate">
                {podcast}
            </div>
            <span>{getTimeAgo(item.pubDate)}</span>
        </div>
        <div class=" card-icons flex space-x-2">
            <button class="btn btn-sm btn-ghost"
                ><i class="fa-regular fa-plus"></i></button
            >
            <button class="btn btn-sm btn-ghost"
                ><i class="fa-regular fa-bookmark"></i></button
            >
        </div>
    </div>
</div>
