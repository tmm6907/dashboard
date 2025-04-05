<script>
    import { alertState } from "$lib/state.svelte";

    // $effect(() => {
    //     if (!alertState.persist) {
    //         console.log("start closing");
    //         setTimeout(() => {
    //             console.log("closing");
    //             alertState.showAlert = false;
    //         }, alertState.duration);
    //     }
    // });
    const closeAlert = () => {
        alertState.showAlert = false;
    };
</script>

<div
    class="relative alert {alertState.type}
    mt-8 mx-auto min-w-[12ch] w-fit
    xl:max-w-1/3 max-w-5/6 z-20
    px-4 py-2 text-wrap {alertState.showAlert ? '' : 'hidden'}"
>
    <span class="mr-12">{alertState.msg}</span>
    {#if alertState.closable}
        <button
            class="btn btn-ghost absolute top-0 right-0"
            onclick={closeAlert}><i class="fa fa-close"></i></button
        >
    {/if}
</div>
