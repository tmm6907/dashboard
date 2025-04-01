<script>
    import { onMount } from "svelte";
    let { item } = $props();
    const share = () => {
        if (navigator.share) {
            navigator
                .share({
                    title: item.title,
                    url: window.location.href,
                })
                .then(() => console.log("Share successful"))
                .catch((error) => console.log("Sharing failed", error));
        } else {
            alert("Sharing is not supported on this browser.");
        }
    };
    onMount(() => {
        let element = document.getElementById("action-btns");
        let timeout;
        window.addEventListener("scroll", () => {
            // Hide the element when scrolling
            element.classList.remove("opacity-100");
            element.classList.add("opacity-0");

            // Clear any previous timeout to track when scrolling stops
            clearTimeout(timeout);

            // Set a new timeout to show the element after scrolling stops
            timeout = setTimeout(() => {
                element.classList.remove("opacity-0");
                element.classList.add("opacity-100");
            }, 100); // Adjust the time as needed
        });
    });
</script>

<div
    id="action-btns"
    class="flex flex-col-reverse gap-4 w-fit fixed bottom-24 right-6 z-20 transition-opacity duration-300 ease-in-out opacity-100"
>
    <button id="shareButton" on:click={share} class="btn btn-info btn-circle">
        <i class="fa fa-share"></i>
    </button>
    <button class="btn btn-secondary btn-circle">
        <i class="fa fa-plus"></i>
    </button>
    <button class="btn btn-warning btn-circle">
        <i class="fa-solid fa-bookmark"></i>
    </button>
</div>
