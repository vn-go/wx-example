export function getAppMenuData() {
    return [
        {
            title: "System",
            children: [
                {
                    id: 1,
                    title: "Users",
                    pathname: "/system/users"
                }, {
                    id: 2,
                    title: "Roles",
                    pathname: "/system/roles"
                }
            ]
        }, {
            title: "Documents",
            children: [
                {
                    id: 1,
                    title: "Persolnal",
                    pathname: "/documents/persolnal"
                }, {
                    id: 2,
                    title: "Public",
                    pathname: "/documents/public"
                }
            ]
        }
    ]
}