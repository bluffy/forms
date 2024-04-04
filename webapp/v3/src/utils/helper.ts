export function propName(prop, value) {
    for (var i in prop) {
        if (typeof prop[i] == 'object') {
            if (propName(prop[i], value)) {
                return res;
            }
        } else {
            if (prop[i] == value) {
                var res = i;
                return res;
            }
        }
    }
    return undefined;
}


export function getPropertyName<T extends object>(o: T, expression: (x : { [Property in keyof T]: string }) => string) {
    const res = {} as { [Property in keyof T]: string };
    Object.keys(o).map(k => res[k as keyof T] = k);
    return expression(res);
  }
  