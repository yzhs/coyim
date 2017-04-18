package definitions

func init() {
	add(`StartVerificationNotification`, &defStartVerificationNotification{})
}

type defStartVerificationNotification struct{}

func (*defStartVerificationNotification) String() string {
	return `<interface>
  <object class="GtkInfoBar" id="infobar">
    <property name="message-type">GTK_MESSAGE_WARNING</property>
    <child internal-child="content_area">
      <object class="GtkBox" id="box">
        <property name="homogeneous">false</property>
        <child>
          <object class="GtkImage">
            <property name="file">build/images/alert.png</property>
            <property name="margin-start">10</property>
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="message">
            <property name="wrap">true</property>
          </object>
        </child>
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkBox" id="button_box">
        <child>
          <object class="GtkButton" id="button_verify">
            <property name="label" translatable="yes">Validate Secure Channel</property>
            <property name="can-default">true</property>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}
