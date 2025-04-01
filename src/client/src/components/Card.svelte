<script>
    import { getTimeAgo } from "$lib";
    import { onMount } from "svelte";
    let { item = {}, feed = true } = $props();
    let boomarked = $state(false);
    const goToFeedItem = () => {
        sessionStorage.setItem("selectedItem", `${item.id}`);
        window.location.href = "/items/" + item.id;
    };
    const goToCollection = () => {
        if (!item.id) {
            window.location.href = "/collections/";
        }
        sessionStorage.setItem("selectedItem", `${item.id}`);
        window.location.href = "/collections/" + item.id;
    };
    const bookmarkFeedItem = async () => {
        console.log("entered");
        if (!item.id) {
            console.error("feed item id not found");
            return;
        }
        try {
            let response = await fetch(
                `http://localhost:8080/api/feeds/items/${item.id}/bookmark`,
                {
                    method: "POST",
                },
            );
            if (response.status == 302) {
                window.location.href = await response.text();
                return;
            }
            if (response.status != 200) {
                throw new Error(await response.text());
            }
            boomarked = !boomarked;
        } catch (e) {
            console.error(e);
        }
    };
    onMount(() => {
        console.log("card item:", item);
        const my_modal_2 = document.getElementById("my_modal_2");
    });
</script>

<!-- Open the modal using ID.showModal() method -->

{#if feed}
    <div
        role="button"
        class="text-base-content min-w-[32ch] max-w-[48ch]"
        onclick={goToFeedItem}
    >
        <img
            src={item.image
                ? item.image
                : "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"}
            alt={item.alt_text ? item.alt_text : "feed item"}
            loading="lazy"
            class="h-48 w-full object-fill"
        />
        <span class=" text-wrap line-clamp-2">{item.title}</span>
        <div class="flex justify-between py-2">
            <div class="flex prose gap-2">
                <div class="max-w-[24ch] text-secondary truncate">
                    {item.feed_name}
                </div>
                |
                <span>{getTimeAgo(item["pub_date"])}</span>
            </div>
            <div class=" card-icons flex space-x-2">
                <button class="btn btn-sm btn-ghost"
                    ><i class="fa-regular fa-plus"></i></button
                >
                <button class="btn btn-sm btn-ghost" onclick={bookmarkFeedItem}
                    ><i id="bookmark-icon" class="fa-regular fa-bookmark"
                    ></i></button
                >
            </div>
        </div>
    </div>
{:else}
    <dialog id="my_modal_2" class="modal">
        <div class="modal-box">
            <h3 class="text-lg font-bold">Hello!</h3>
            <p class="py-4">Press ESC key or click outside to close</p>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>
    <button
        class="flex justify-center items-center h-32 border-0 rounded-lg bg-base-300 pb-4 min-w-[32ch] max-w-[48ch]"
        onclick={item.id ? goToCollection : my_modal_2.showModal}
    >
        <span
            class="prose text-4xl font-bold text-neutral-400 focus:underline-offset-4 focus:underline-4 focus:decoration-secondary hover:underline-offset-4 hover:underline hover:decoration-secondary"
        >
            {#if item.name}
                {item.name}
            {:else}
                <i class="fa fa-plus"></i>
            {/if}
        </span>
    </button>
{/if}
