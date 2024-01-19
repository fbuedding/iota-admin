// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package template

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func InternalAttributes() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<internal-attributes-input></internal-attributes-input>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func InternalAttributeTemplate() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<template id=\"internal-attributes-container-template\"><div class=\"pl-5 pt-5\"><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := `Internal Attributes:`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"internal-attributes\" class=\"\"></div><div class=\"pt-5\"><input id=\"add-btn\" type=\"button\" value=\"Add\" class=\"btn\"> <input id=\"rm-btn\" type=\"button\" value=\"Remove\" class=\"btn\" disabled></div></div></template><template id=\"internal-attribute-template\"><div class=\"collapse collapse-arrow mt-5 border border-neutral-500/40\"><input type=\"checkbox\" checked><div id=\"static-attribute-title\" class=\"collapse-title text-xl font-medium\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var4 := `%index`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"collapse-content\"><label class=\"form-control\"><div class=\"label\"><span class=\"label-text\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var5 := `Free JSON format`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var5)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div><textarea name=\"internal_attributes[%index]\" class=\"textarea textarea-bordered h-24\" placeholder=\"JSON\" required></textarea></label></div></div></template>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InternalAttributesScript().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func InternalAttributesScript() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_InternalAttributesScript_ddf7`,
		Function: `function __templ_InternalAttributesScript_ddf7(){// Create a class for the element
  class InternalAttribute extends HTMLElement {
    static formAssociated = true;

    constructor() {
      // Always call super first in constructor
      super();
      this.internals_ = this.attachInternals();
      this.state = {
        count: 0,
      };
    }

    connectedCallback() {
      this.attributeTemplate = document.getElementById("internal-attribute-template");
      let template = document.getElementById("internal-attributes-container-template");
      this.appendChild(template.content.cloneNode(true));
      this.attributesContainer = this.querySelector("#internal-attributes");

      this.addBtn = this.querySelector("#add-btn");
      this.addBtn.onclick = () => { this.addAttribute() };
      this.removeBtn = this.querySelector("#rm-btn");
      this.removeBtn.onclick = () => { this.removeGenericAttribute() };
    }

    disconnectedCallback() {}

    adoptedCallback() {}

    attributeChangedCallback(name, oldValue, newValue) {}
    addAttribute() {
      const el = this.attributeTemplate.content.cloneNode(true).firstElementChild;
      el.innerHTML = el.innerHTML.replaceAll("%index", this.state.count);
      this.attributesContainer.appendChild(el);
      el.scrollIntoView({behavior: "smooth", block:"center"});
      
      this.state.count++;
      this.removeBtn.disabled = (this.state.count < 1);
    }

    removeGenericAttribute() {
      this.attributesContainer.removeChild(this.attributesContainer.lastChild)
      this.state.count--;
      this.removeBtn.disabled = (this.state.count < 1);
    }
  }

  customElements.define("internal-attributes-input", InternalAttribute);}`,
		Call:       templ.SafeScript(`__templ_InternalAttributesScript_ddf7`),
		CallInline: templ.SafeScriptInline(`__templ_InternalAttributesScript_ddf7`),
	}
}
