package definitions

func init() {
	add(`EnterPIN`, &defEnterPIN{})
}

type defEnterPIN struct{}

func (*defEnterPIN) String() string {
	return `<interface>
  <object class="GtkEntryBuffer" id="pin_entry_buffer">
    <signal name="inserted-text" handler="text_changing"/>
    <signal name="deleted-text" handler="text_changing"/>
  </object>
  <object class="GtkDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="box">
        <property name="border-width">10</property>
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkBox" id="notification-area">
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object  class="GtkImage">
            <property name="file">build/images/smp.svg</property>
            <property name="margin-bottom">10</property>
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="verification_message">
            <property name="label" translatable="yes"/>
          </object>
        </child>
        <child>
          <object class="GtkLabel">
            <property name="label" translatable="yes">Isso pode ser usado somente uma vez.</property>
            <property name="margin-bottom">10</property>
          </object>
        </child>
        <child>
          <object class="GtkGrid">
            <property name="halign">GTK_ALIGN_CENTER</property>
            <child>
              <object  class="GtkLabel">
                <property name="label" translatable="yes">CÓDIGO:</property>
                <property name="justify">GTK_JUSTIFY_RIGHT</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">0</property>
                <property name="width">2</property>
              </packing>
            </child>
            <child>
              <object  class="GtkEntry" id="pin">
                <property name="buffer">pin_entry_buffer</property>
                <property name="margin-bottom">10</property>
              </object>
              <packing>
                <property name="top-attach">0</property>
              </packing>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkLabel">
            <property name="label" translatable="yes">Seu contato deverá ter compartilhado isso com você anteriormente. Se não, tente:</property>
            <property name="margin-bottom">5</property>
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
                    <property name="label" translatable="yes">Perguntar pessoalmente</property>
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
                    <property name="label" translatable="yes">Perguntar em outra aplicação criptografada</property>
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
                    <property name="label" translatable="yes">Perguntar pelo telefone</property>
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
        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="button_submit">
                <property name="label" translatable="yes">Verificar o código</property>
                <signal name="clicked" handler="close_share_pin"/>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
    <action-widgets>
      <action-widget response="ok">button_submit</action-widget>
    </action-widgets>
  </object>
</interface>
`
}
