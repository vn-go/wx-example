<template>
    <div class="flex flex-col h-full">
        <div class="pl-2 pr-2">Header</div>
        <div class="flex-1 pl-2 pr-2">
             <slot name="main">main</slot>

        </div>
        
        <div class="pl-2 pr-2">Footer</div>
    </div>
  </template>
  