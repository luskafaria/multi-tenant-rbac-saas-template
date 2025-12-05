'use client'

import * as SheetPrimitive from '@radix-ui/react-dialog'
import { X } from 'lucide-react'
import { useRouter } from 'next/navigation'
import * as React from 'react'

import { cn } from '@/lib/utils'

import {
  type SheetContentProps,
  SheetOverlay,
  SheetPortal,
  sheetVariants,
} from './ui/sheet'

export const InterceptedSheetContent = React.forwardRef<
  React.ElementRef<typeof SheetPrimitive.Content>,
  SheetContentProps
>(({ side = 'right', className, children, ...props }, ref) => {
  const router = useRouter()

  function onDismiss() {
    router.back()
  }
  return (
    <SheetPortal>
      <SheetOverlay />
      <SheetPrimitive.Content
        onEscapeKeyDown={onDismiss}
        onPointerDownOutside={onDismiss}
        ref={ref}
        className={cn(sheetVariants({ side }), className)}
        {...props}
      >
        {children}
        <SheetPrimitive.Close
          onClick={onDismiss}
          className="ring-offset-background focus:ring-ring data-[state=open]:bg-secondary absolute top-4 right-4 rounded-sm opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none"
        >
          <X className="h-4 w-4" />
          <span className="sr-only">Close</span>
        </SheetPrimitive.Close>
      </SheetPrimitive.Content>
    </SheetPortal>
  )
})
InterceptedSheetContent.displayName = SheetPrimitive.Content.displayName
