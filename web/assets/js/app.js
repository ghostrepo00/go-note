const app = (function () {    
    const bindingDialogButtons = () => {
        const dialog = document.querySelector("dialog");
        const showButton = document.getElementById("deleteButton")
        const cancelButton = document.getElementById("cancelButton")
        
        if (!!showButton) {
            showButton.addEventListener("click", () => {
                dialog.showModal();
            });
        }

        cancelButton.addEventListener("click", () => {
            dialog.close();
        });
    }

    htmx.onLoad(function(elt) {
        if (htmx.findAll(elt, "#contentRoot").length > 0) {
            bindingDialogButtons();
        }
    })

    const markedOptions = {
        pedantic: false,
        gfm: true,
        breaks: true
    };

    return {
        closeDialog: () => {dialog.close();},
        markdownMode: (element) => {
            element.innerHTML = marked.parse([].slice.call(element.childNodes).map(a => a.textContent).join('\n'), markedOptions);
            document.getElementById("xData")._x_dataStack[0].isEdit = false;
        },
        editMode: (element, content) => {
            element.innerHTML = content
            document.getElementById("xData")._x_dataStack[0].isEdit = true;
        }
    }
})();