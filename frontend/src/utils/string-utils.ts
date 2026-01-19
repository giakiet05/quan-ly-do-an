/**
 * Chuyển đổi một chuỗi từ snake_case sang camelCase
 */
export const toCamel = (str: string): string => {
    return str.replace(/([-_][a-z])/ig, ($1) => {
        return $1.toUpperCase()
            .replace('-', '')
            .replace('_', '');
    });
};

/**
 * Chuyển đổi toàn bộ Object hoặc Array từ snake_case sang camelCase (Đệ quy)
 */
export const keysToCamel = (obj: any): any => {
    if (Array.isArray(obj)) {
        return obj.map((v) => keysToCamel(v));
    } else if (obj !== null && obj !== undefined && obj.constructor === Object) {
        return Object.keys(obj).reduce(
            (result, key) => ({
                ...result,
                [toCamel(key)]: keysToCamel(obj[key]),
            }),
            {},
        );
    }
    return obj;
};