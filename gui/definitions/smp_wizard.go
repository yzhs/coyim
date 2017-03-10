package definitions

func init() {
	add(`SMPWizard`, &defSMPWizard{})
}

type defSMPWizard struct{}

func (*defSMPWizard) String() string {
	return `<interface>
   <object class="GtkAssistant" id="SMPWizard">
    <signal name="close" handler="on_close_signal" />
    <signal name="cancel" handler="on_cancel_signal" />
    <signal name="escape" handler="on_escape_signal" />
    <child>
        <object class="GtkBox">
            <property name="homogeneous">false</property>
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
            <property name="spacing">18</property>
            <child>
                <object class="GtkImage">
                    <property name="file">build/images/maninthemiddle.png</property>
                </object>
                <packing>
                    <property name="padding">60</property>
                </packing>
            </child>
            <child>
                <object class="GtkLabel">
                    <property name="visible">True</property>
                    <property name="label" translatable="yes">Make sure there is no one else reading your messages.</property>
                </object>
            </child>
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
            <object class="GtkImage">
                <property name="file">build/images/Pin_1.png</property>
            </object>
        </child>
        <child>
            <object class="GtkLabel">
                <property name="visible">True</property>
                <property name="label" translatable="yes">Share this pin only with your contact.</property>
            </object>
        </child>
        <child>
            <object class="GtkLabel">
                <property name="visible">True</property>
                <property name="label" translatable="yes">It can only be used once.</property>
            </object>
        </child>
        <child>
            <object class="GtkLabel" id="PinLabel">
                <property name="visible">True</property>
                <property name="selectable">True</property>
                <property name="label">placeholder</property>
            </object>
        </child>
        <child>
            <object class="GtkGrid">
            <property name="column-spacing">6</property>
            <property name="row-spacing">2</property>
            <child>
                <object class="GtkImage">
                    <property name="file">build/images/padlock.png</property>
                </object>
                <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                </packing>
            </child>
            <child>
                <object class="GtkLabel">
                    <property name="visible">True</property>
                    <property name="label" translatable="yes">Share in person</property>
                    <property name="justify">GTK_JUSTIFY_LEFT</property>
                    <property name="halign">GTK_ALIGN_START</property>
                </object>
                <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                </packing>
            </child>
            <child>
                <object class="GtkImage">
                    <property name="file">build/images/padlock.png</property>
                </object>
                <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                </packing>
            </child>
            <child>
                <object class="GtkLabel">
                    <property name="visible">True</property>
                    <property name="label" translatable="yes">Share through an encrypted channel</property>
                    <property name="justify">GTK_JUSTIFY_LEFT</property>
                    <property name="halign">GTK_ALIGN_START</property>
                </object>
                <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                </packing>
            </child>
            <child>
                <object class="GtkImage">
                    <property name="file">build/images/alert.png</property>
                </object>
                <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                </packing>
            </child>
            <child>
                <object class="GtkLabel">
                    <property name="visible">True</property>
                    <property name="label" translatable="yes">Share in a phone call</property>
                    <property name="justify">GTK_JUSTIFY_LEFT</property>
                    <property name="halign">GTK_ALIGN_START</property>
                </object>
                <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">2</property>
                </packing>
            </child>
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
