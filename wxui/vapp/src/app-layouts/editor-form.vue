<template>
    <div class="flex flex-col w-full h-full">
        <div>
            <slot name="header">
                Header
            </slot>
        </div>
        <div class="h-full flex flex-col">
            <slot name="body">
                Body
            </slot>

        </div>
        <div>
            <slot name="footer">
                Footer
            </slot>
        </div>

    </div>
</template>