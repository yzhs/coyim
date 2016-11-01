package definitions

func init() {
	add(`AddContact2`, &defAddContact2{})
}

type defAddContact2 struct{}

func (*defAddContact2) String() string {
	return `<interface>
   <object class="GtkAssistant" id="AddContact">
    <signal name="close" handler="on_close_signal" />
    <signal name="cancel" handler="on_cancel_signal" />
    <signal name="escape" handler="on_escape_signal" />
    <child>
        <object class="GtkLabel" id="label1">
            <property name="visible">True</property>
            <property name="label" translatable="yes">Make sure there is no one else reading your messages.</property>
        </object>
        <packing>
            <property name="page_type">content</property>
        </packing>
    </child>
    <child>
        <object class="GtkBox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <property name="spacing">6</property>
        <child>
            <object class="GtkLabel" id="share">
                <property name="visible">True</property>
                <property name="label" translatable="yes">Share this pin only with your contact</property>
            </object>
        </child>
        <child>
            <object class="GtkLabel" id="onlyOnce">
                <property name="visible">True</property>
                <property name="label" translatable="yes">It can only be used once.</property>
            </object>
        </child>
        </object>
        <packing>
            <property name="page_type">summary</property>
        </packing>
    </child>
  </object>
</interface>
`
}
