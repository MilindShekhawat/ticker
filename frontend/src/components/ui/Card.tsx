import { ComponentChildren, HTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const cardVariants = cva(`bg-zinc-900 text-zinc-100 border border-zinc-800 transition`, {
    variants: {
        interactive: {
            true: `cursor-pointer hover:border-lime-400`,
            false: ``,
        },
    },
    defaultVariants: {
        interactive: false,
    },
});

type CardProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
} & VariantProps<typeof cardVariants>;

type CardHeaderProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
};

type CardTitleProps = HTMLAttributes<HTMLHeadingElement> & {
    children: ComponentChildren;
};

type CardDescriptionProps = HTMLAttributes<HTMLParagraphElement> & {
    children: ComponentChildren;
};

type CardContentProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
};

type CardFooterProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
};
export function Card({ children, class: className, onClick, ...props }: CardProps) {
    return (
        <div
            {...props}
            onClick={onClick}
            class={cardVariants({
                interactive: Boolean(onClick),
                class: className,
            })}>
            {children}
        </div>
    );
}

export function CardHeader({ children, class: className, ...props }: CardHeaderProps) {
    return (
        <div
            {...props}
            class={["flex flex-col space-y-2 p-6", className].filter(Boolean).join(" ")}>
            {children}
        </div>
    );
}

export function CardTitle({ children, class: className, ...props }: CardTitleProps) {
    return (
        <h3
            {...props}
            class={["text-2xl font-semibold leading-none", className].filter(Boolean).join(" ")}>
            {children}
        </h3>
    );
}

export function CardDescription({ children, class: className, ...props }: CardDescriptionProps) {
    return (
        <p {...props} class={["text-sm text-zinc-400", className].filter(Boolean).join(" ")}>
            {children}
        </p>
    );
}

export function CardContent({ children, class: className, ...props }: CardContentProps) {
    return (
        <div
            {...props}
            class={["p-6 pt-0 text-sm text-zinc-400", className].filter(Boolean).join(" ")}>
            {children}
        </div>
    );
}

export function CardFooter({ children, class: className, ...props }: CardFooterProps) {
    return (
        <div {...props} class={["flex items-center p-6 pt-0", className].filter(Boolean).join(" ")}>
            {children}
        </div>
    );
}
