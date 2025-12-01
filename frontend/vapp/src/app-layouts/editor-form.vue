<template>
    <div class="flex flex-col w-full h-full">
        <div class="widget-container">
            <slot name="header">
                Header
            </slot>
        </div>
        <div class="h-full flex flex-col m-0 p-0">
            
            <slot name="body">
                Body
            </slot>

        </div>
        <div class="w-full  m-0 p-0">
          
            <slot name="footer">
                Footer
            </slot>
        </div>

    </div>
</template>