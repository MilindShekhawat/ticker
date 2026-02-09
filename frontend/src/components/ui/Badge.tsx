import { ComponentChildren, HTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const badgeVariants = cva(
    `inline-flex items-center border-transparent px-1.5 py-0.5 text-xs font-bold transition-colors`,
    {
        variants: {
            variant: {
                primary: `bg-lime-400/40 text-lime-400`,
                destructive: `bg-orange-600/40 text-orange-500`,
                success: `bg-green-500/40 text-green-400`,
                warning: `bg-yellow-500/40 text-yellow-400`,
                blue: `bg-sky-500/40 text-sky-400`,
                purple: `bg-purple-500/40 text-purple-400`,
            },
        },
        defaultVariants: {
            variant: `primary`,
        },
    },
);

type BadgeProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
} & VariantProps<typeof badgeVariants>;

export function Badge({ children, variant, class: className, ...props }: BadgeProps) {
    return (
        <div {...props} class={badgeVariants({ variant, class: className })}>
            {children}
        </div>
    );
}
