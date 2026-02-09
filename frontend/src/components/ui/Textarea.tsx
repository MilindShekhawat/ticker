import { TextareaHTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const textareaVariants = cva(
    `inline-flex w-full min-h-24 px-4 py-2 text-sm mt-1 font-medium border bg-zinc-900 transition outline-none resize-y disabled:pointer-events-none disabled:opacity-40`,
    {
        variants: {
            state: {
                default: `placeholder:text-zinc-500 text-zinc-100 border-zinc-700 focus:border-zinc-300 focus:text-zinc-300`,
                error: `placeholder:text-orange-800 text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`,
            },
        },
        defaultVariants: {
            state: `default`,
        },
    },
);

type TextareaProps = TextareaHTMLAttributes & {
    error?: string;
} & VariantProps<typeof textareaVariants>;

export function Textarea({ class: className, error, state, ...props }: TextareaProps) {
    return (
        <div class="w-full">
            <textarea
                {...props}
                class={textareaVariants({
                    state: error ? "error" : state,
                    class: className,
                })}
            />
            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
