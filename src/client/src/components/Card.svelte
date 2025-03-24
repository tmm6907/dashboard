<script>
    let { podcast, item } = $props();

    const getTimeAgo = (time) => {
        let date = new Date(time.replace(" ", "T") + "Z");

        const now = new Date();
        const diff = Math.floor((now - date) / 1000); // Difference in seconds

        const days = Math.floor(diff / (60 * 60 * 24));
        const hours = Math.floor(diff / (60 * 60));
        const minutes = Math.floor(diff / 60);

        if (days > 0) return `${days}d ago`;
        if (hours > 0) return `${hours}h ago`;
        console.log(time);
        return `${minutes}m ago`;
    };
</script>

<div class="flex flex-col text-base-content pb-4 max-w-[48ch]">
    <img
        src={item.image
            ? item.image
            : "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
        alt={item.altText ? item.altText : "feed item"}
        loading="lazy"
        class="h-48 w-full"
        style=" margin-inline: auto; object-fit:cover;"
    />
    <span>{item.title}</span>
    <div class="prose">
        <span class="text-secondary">{podcast}</span> |
        <span>{getTimeAgo(item.pubDate)}</span>
    </div>
</div>
