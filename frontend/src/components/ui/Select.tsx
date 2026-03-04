import { SelectHTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const selectVariants = cva(
    `inline-flex w-full h-8 px-3 text-sm font-medium border bg-zinc-900 transition outline-none disabled:pointer-events-none disabled:opacity-40`,
    {
        variants: {
            state: {
                default: `text-zinc-100 border-zinc-700 focus:border-lime-400 focus:text-lime-400`,
                error: `text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`,
            },
        },
        defaultVariants: {
            state: `default`,
        },
    },
);

type SelectProps = SelectHTMLAttributes<HTMLSelectElement> & {
    error?: string;
} & VariantProps<typeof selectVariants>;

export function Select({ class: className, error, state, children, ...props }: SelectProps) {
    return (
        <div class="w-full">
            <select
                {...props}
                class={selectVariants({
                    state: error ? "error" : state,
                    class: className,
                })}>
                {children}
            </select>
            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
