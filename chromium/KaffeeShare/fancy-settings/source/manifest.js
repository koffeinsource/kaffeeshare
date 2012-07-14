// SAMPLE
this.manifest = {
    "name": "KaffeeShare",
    "icon": "../../comic_30x30.png",
    "settings": [
        {
            "tab": i18n.get("sharing"),
            "group": i18n.get("backend"),
            "name": "url",
            "type": "text",
            "label": i18n.get("url_label"),
            "text": ""
        },
        {
            "tab": i18n.get("sharing"),
            "group": i18n.get("backend"),
            "name": "url_description",
            "type": "description",
            "text": i18n.get("url_text")
        }
    ],
    "alignment": [
        [
            "url"
        ]
    ]
};
