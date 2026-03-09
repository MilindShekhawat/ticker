import { InputHTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const inputVariants = cva(
    `inline-flex w-full h-8 px-4 text-sm font-medium border bg-zinc-900 transition outline-none disabled:pointer-events-none disabled:opacity-40`,
    {
        variants: {
            state: {
                default: `placeholder:text-zinc-500 text-zinc-100 border-zinc-700 focus:border-lime-400 focus:text-lime-400`,
                error: `placeholder:text-orange-800 text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`,
            },
        },
        defaultVariants: {
            state: `default`,
        },
    },
);

type InputProps = InputHTMLAttributes & {
    error?: string;
} & VariantProps<typeof inputVariants>;

export function Input({ class: className, error, state, ...props }: InputProps) {
    return (
        <div class="w-full">
            <input
                {...props}
                class={inputVariants({
                    state: error ? "error" : state,
                    class: className,
                })}
            />
            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
