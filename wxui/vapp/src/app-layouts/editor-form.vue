<template>
    <div class="flex flex-col w-full h-full">
        <div>
            <slot name="header">
                Header
            </slot>
        </div>
        <div class="h-full flex flex-col m-0 p-0">
            
            <slot name="body">
                Body
            </slot>

        </div>
        <div class="w-full  m-0 p-0 flex flex-row">
            <div class="w-full"></div>
            <slot name="footer">
                Footer
            </slot>
        </div>

    </div>
</template>