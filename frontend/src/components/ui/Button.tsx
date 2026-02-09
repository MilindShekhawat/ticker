import { ComponentChildren } from "preact";
import { ButtonHTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const buttonVariants = cva(
    `inline-flex uppercase items-center justify-center font-medium border transition disabled:pointer-events-none disabled:opacity-40`,
    {
        variants: {
            variant: {
                primary: `bg-lime-400 text-zinc-950 border-lime-400 hover:border-lime-400 hover:text-lime-400 hover:bg-zinc-950`,
                white: `bg-zinc-100 text-zinc-950 border-zinc-100 hover:border-lime-400 hover:text-lime-400 hover:bg-zinc-950`,
                outline: `bg-zinc-950 text-zinc-100 border-zinc-100 hover:border-lime-400 hover:text-lime-400`,
                black: `bg-zinc-950 text-zinc-100 border-zinc-950 hover:border-lime-400 hover:text-lime-400`,
                invert: `bg-zinc-950 text-zinc-100 border-zinc-950 hover:border-lime-400 hover:text-zinc-950 hover:bg-lime-400`,
                destructive: `bg-orange-600 text-zinc-100 border-orange-600 hover:bg-zinc-950 hover:border-orange-600 hover:text-orange-600`,
            },
            size: {
                xs: `h-4 px-0 text-sm`,
                sm: `h-6 px-2 text-sm`,
                md: `h-8 px-4 text-sm`,
                lg: `h-10 px-8 text-base`,
                xl: `h-12 px-12 text-base`,
            },
        },
        defaultVariants: {
            variant: `primary`,
            size: `md`,
        },
    },
);

type ButtonProps = ButtonHTMLAttributes & {
    children: ComponentChildren;
} & VariantProps<typeof buttonVariants>;

export function Button({ children, variant, size, class: className, ...props }: ButtonProps) {
    return (
        <button {...props} class={buttonVariants({ variant, size, class: className })}>
            {children}
        </button>
    );
}
